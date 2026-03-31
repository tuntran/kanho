import { create } from "zustand";
import type { Board, Card, Column, WSEvent } from "@/types/board";

interface OpSnapshot {
  cardId: string;
  prevColumnId: string;
  prevPosition: string;
  prevCardsByColumn: Record<string, string[]>;
}

interface BoardState {
  boardId: string | null;
  board: Board | null;
  columns: Column[];
  cards: Record<string, Card>;
  cardsByColumn: Record<string, string[]>;
  pendingOps: Map<string, OpSnapshot>;
  processedEventIds: string[];
  wsConnected: boolean;
}

interface BoardActions {
  initBoard: (board: Board, columns: Column[], cards: Card[]) => void;
  reset: () => void;
  moveCardOptimistic: (
    cardId: string,
    toColumnId: string,
    afterCardId: string | null,
    beforeCardId: string | null,
  ) => string;
  confirmOp: (opId: string) => void;
  rollbackOp: (opId: string) => void;
  addCard: (card: Card) => void;
  updateCard: (card: Card) => void;
  removeCard: (cardId: string) => void;
  addColumn: (column: Column) => void;
  handleServerEvent: (event: WSEvent) => void;
  setWsConnected: (v: boolean) => void;
}

const BUFFER_MAX = 200;

function sortedCardIds(
  cards: Record<string, Card>,
  ids: string[],
): string[] {
  return [...ids].sort((a, b) => {
    const ca = cards[a];
    const cb = cards[b];
    if (!ca || !cb) return 0;
    return ca.position.localeCompare(cb.position);
  });
}

function addToBuffer(buf: string[], id: string): string[] {
  const next = [...buf, id];
  if (next.length > BUFFER_MAX) {
    return next.slice(next.length - BUFFER_MAX);
  }
  return next;
}

let opCounter = 0;

export const useBoardStore = create<BoardState & BoardActions>((set, get) => ({
  boardId: null,
  board: null,
  columns: [],
  cards: {},
  cardsByColumn: {},
  pendingOps: new Map(),
  processedEventIds: [],
  wsConnected: false,

  initBoard: (board, columns, cards) => {
    const cardsMap: Record<string, Card> = {};
    const byColumn: Record<string, string[]> = {};

    for (const col of columns) {
      byColumn[col.id] = [];
    }

    for (const card of cards) {
      if (card.archived_at) continue;
      cardsMap[card.id] = card;
      if (!byColumn[card.column_id]) {
        byColumn[card.column_id] = [];
      }
      byColumn[card.column_id]!.push(card.id);
    }

    // Sort each column's cards by position
    for (const colId of Object.keys(byColumn)) {
      byColumn[colId] = sortedCardIds(cardsMap, byColumn[colId]!);
    }

    const sortedColumns = [...columns].sort((a, b) =>
      a.position.localeCompare(b.position),
    );

    set({
      boardId: board.id,
      board,
      columns: sortedColumns,
      cards: cardsMap,
      cardsByColumn: byColumn,
      pendingOps: new Map(),
      processedEventIds: [],
    });
  },

  reset: () =>
    set({
      boardId: null,
      board: null,
      columns: [],
      cards: {},
      cardsByColumn: {},
      pendingOps: new Map(),
      processedEventIds: [],
      wsConnected: false,
    }),

  moveCardOptimistic: (cardId, toColumnId, afterCardId, beforeCardId) => {
    const state = get();
    const card = state.cards[cardId];
    if (!card) return "";

    const opId = `op-${++opCounter}`;
    const fromColumnId = card.column_id;

    // Snapshot for rollback
    const snapshot: OpSnapshot = {
      cardId,
      prevColumnId: fromColumnId,
      prevPosition: card.position,
      prevCardsByColumn: {
        [fromColumnId]: [...(state.cardsByColumn[fromColumnId] ?? [])],
        [toColumnId]: [...(state.cardsByColumn[toColumnId] ?? [])],
      },
    };

    // Remove from old column
    const fromList = (state.cardsByColumn[fromColumnId] ?? []).filter(
      (id) => id !== cardId,
    );

    // Insert into new column at correct position
    let toList =
      fromColumnId === toColumnId
        ? fromList
        : [...(state.cardsByColumn[toColumnId] ?? [])].filter(
            (id) => id !== cardId,
          );

    let insertIdx = toList.length; // default: end
    if (afterCardId && beforeCardId) {
      const afterIdx = toList.indexOf(afterCardId);
      insertIdx = afterIdx >= 0 ? afterIdx + 1 : toList.length;
    } else if (afterCardId) {
      const afterIdx = toList.indexOf(afterCardId);
      insertIdx = afterIdx >= 0 ? afterIdx + 1 : toList.length;
    } else if (beforeCardId) {
      const beforeIdx = toList.indexOf(beforeCardId);
      insertIdx = beforeIdx >= 0 ? beforeIdx : 0;
    }

    toList = [...toList.slice(0, insertIdx), cardId, ...toList.slice(insertIdx)];

    const newCardsByColumn = { ...state.cardsByColumn };
    newCardsByColumn[fromColumnId] = fromList;
    newCardsByColumn[toColumnId] = toList;

    const newCards = {
      ...state.cards,
      [cardId]: { ...card, column_id: toColumnId },
    };

    const newPending = new Map(state.pendingOps);
    newPending.set(opId, snapshot);

    set({
      cards: newCards,
      cardsByColumn: newCardsByColumn,
      pendingOps: newPending,
    });

    return opId;
  },

  confirmOp: (opId) => {
    const newPending = new Map(get().pendingOps);
    newPending.delete(opId);
    set({ pendingOps: newPending });
  },

  rollbackOp: (opId) => {
    const state = get();
    const snapshot = state.pendingOps.get(opId);
    if (!snapshot) return;

    const newCards = { ...state.cards };
    const card = newCards[snapshot.cardId];
    if (card) {
      newCards[snapshot.cardId] = {
        ...card,
        column_id: snapshot.prevColumnId,
        position: snapshot.prevPosition,
      };
    }

    const newCardsByColumn = { ...state.cardsByColumn };
    for (const [colId, ids] of Object.entries(snapshot.prevCardsByColumn)) {
      newCardsByColumn[colId] = ids;
    }

    const newPending = new Map(state.pendingOps);
    newPending.delete(opId);

    set({
      cards: newCards,
      cardsByColumn: newCardsByColumn,
      pendingOps: newPending,
    });
  },

  addCard: (card) => {
    if (card.archived_at) return;
    const state = get();
    const newCards = { ...state.cards, [card.id]: card };
    const colList = [...(state.cardsByColumn[card.column_id] ?? [])];
    if (!colList.includes(card.id)) {
      colList.push(card.id);
    }
    const sorted = sortedCardIds(newCards, colList);
    set({
      cards: newCards,
      cardsByColumn: { ...state.cardsByColumn, [card.column_id]: sorted },
    });
  },

  updateCard: (card) => {
    const state = get();
    set({ cards: { ...state.cards, [card.id]: card } });
  },

  removeCard: (cardId) => {
    const state = get();
    const card = state.cards[cardId];
    if (!card) return;
    const newCards = { ...state.cards };
    delete newCards[cardId];
    const colList = (state.cardsByColumn[card.column_id] ?? []).filter(
      (id) => id !== cardId,
    );
    set({
      cards: newCards,
      cardsByColumn: { ...state.cardsByColumn, [card.column_id]: colList },
    });
  },

  addColumn: (column) => {
    const state = get();
    const newColumns = [...state.columns, column].sort((a, b) =>
      a.position.localeCompare(b.position),
    );
    set({
      columns: newColumns,
      cardsByColumn: { ...state.cardsByColumn, [column.id]: [] },
    });
  },

  handleServerEvent: (event) => {
    const state = get();
    const eventId = event.id ?? event.target_id;
    if (state.processedEventIds.includes(eventId)) return;

    set({ processedEventIds: addToBuffer(state.processedEventIds, eventId) });

    // Events from server are authoritative — but skip if we have pending op for same card
    // to avoid flicker. The pending op will be confirmed/rolled back separately.
  },

  setWsConnected: (v) => set({ wsConnected: v }),
}));
