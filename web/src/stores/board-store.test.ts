import { describe, it, expect, beforeEach } from "vitest";
import { useBoardStore } from "./board-store";
import type { Board, Card, Column } from "@/types/board";

const board: Board = {
  id: "b1",
  project_id: "p1",
  name: "Test Board",
  slug: "test-board",
  visibility: "private",
  created_at: "2025-01-01T00:00:00Z",
  updated_at: "2025-01-01T00:00:00Z",
};

const columns: Column[] = [
  { id: "c1", board_id: "b1", name: "Todo", position: "a", is_done: false, created_at: "2025-01-01T00:00:00Z" },
  { id: "c2", board_id: "b1", name: "Done", position: "b", is_done: true, created_at: "2025-01-01T00:00:00Z" },
];

function makeCard(overrides: Partial<Card> = {}): Card {
  return {
    id: "card1",
    project_id: "p1",
    board_id: "b1",
    column_id: "c1",
    card_number: 1,
    title: "Test Card",
    priority: "none",
    position: "a",
    created_by: "u1",
    created_at: "2025-01-01T00:00:00Z",
    updated_at: "2025-01-01T00:00:00Z",
    ...overrides,
  };
}

describe("boardStore", () => {
  beforeEach(() => {
    useBoardStore.getState().reset();
  });

  it("initializes board with columns and cards sorted by position", () => {
    const cards = [
      makeCard({ id: "card2", position: "b" }),
      makeCard({ id: "card1", position: "a" }),
    ];
    useBoardStore.getState().initBoard(board, columns, cards);

    const state = useBoardStore.getState();
    expect(state.boardId).toBe("b1");
    expect(state.columns).toHaveLength(2);
    expect(state.cardsByColumn["c1"]).toEqual(["card1", "card2"]);
  });

  it("skips archived cards during init", () => {
    const cards = [
      makeCard({ id: "card1" }),
      makeCard({ id: "card2", archived_at: "2025-01-02T00:00:00Z" }),
    ];
    useBoardStore.getState().initBoard(board, columns, cards);

    const state = useBoardStore.getState();
    expect(Object.keys(state.cards)).toEqual(["card1"]);
  });

  it("moves card optimistically between columns", () => {
    const cards = [
      makeCard({ id: "card1", column_id: "c1", position: "a" }),
    ];
    useBoardStore.getState().initBoard(board, columns, cards);

    const opId = useBoardStore.getState().moveCardOptimistic("card1", "c2", null, null);
    const state = useBoardStore.getState();

    expect(opId).toBeTruthy();
    expect(state.cards["card1"]!.column_id).toBe("c2");
    expect(state.cardsByColumn["c1"]).toEqual([]);
    expect(state.cardsByColumn["c2"]).toContain("card1");
  });

  it("confirms op removes it from pending", () => {
    const cards = [makeCard({ id: "card1" })];
    useBoardStore.getState().initBoard(board, columns, cards);

    const opId = useBoardStore.getState().moveCardOptimistic("card1", "c2", null, null);
    expect(useBoardStore.getState().pendingOps.size).toBe(1);

    useBoardStore.getState().confirmOp(opId);
    expect(useBoardStore.getState().pendingOps.size).toBe(0);
  });

  it("rolls back on server error", () => {
    const cards = [makeCard({ id: "card1", column_id: "c1", position: "a" })];
    useBoardStore.getState().initBoard(board, columns, cards);

    const opId = useBoardStore.getState().moveCardOptimistic("card1", "c2", null, null);
    useBoardStore.getState().rollbackOp(opId);

    const state = useBoardStore.getState();
    expect(state.cards["card1"]!.column_id).toBe("c1");
    expect(state.cardsByColumn["c1"]).toContain("card1");
    expect(state.cardsByColumn["c2"]).not.toContain("card1");
    expect(state.pendingOps.size).toBe(0);
  });

  it("deduplicates by eventId in handleServerEvent", () => {
    const cards = [makeCard({ id: "card1" })];
    useBoardStore.getState().initBoard(board, columns, cards);

    const event = {
      board_id: "b1",
      id: "evt-1",
      action: "card.moved",
      target_type: "card",
      target_id: "card1",
    };

    useBoardStore.getState().handleServerEvent(event);
    useBoardStore.getState().handleServerEvent(event);

    const state = useBoardStore.getState();
    // eventId should only appear once in processedEventIds
    const count = state.processedEventIds.filter((id) => id === "evt-1").length;
    expect(count).toBe(1);
  });

  it("adds card to the correct column", () => {
    useBoardStore.getState().initBoard(board, columns, []);

    const newCard = makeCard({ id: "card-new", column_id: "c2", position: "a" });
    useBoardStore.getState().addCard(newCard);

    const state = useBoardStore.getState();
    expect(state.cards["card-new"]).toBeDefined();
    expect(state.cardsByColumn["c2"]).toContain("card-new");
  });

  it("removes card from store and column", () => {
    const cards = [makeCard({ id: "card1", column_id: "c1" })];
    useBoardStore.getState().initBoard(board, columns, cards);

    useBoardStore.getState().removeCard("card1");

    const state = useBoardStore.getState();
    expect(state.cards["card1"]).toBeUndefined();
    expect(state.cardsByColumn["c1"]).not.toContain("card1");
  });

  it("adds column sorted by position", () => {
    useBoardStore.getState().initBoard(board, columns, []);

    const newCol: Column = {
      id: "c1.5",
      board_id: "b1",
      name: "In Progress",
      position: "a1",
      is_done: false,
      created_at: "2025-01-01T00:00:00Z",
    };
    useBoardStore.getState().addColumn(newCol);

    const state = useBoardStore.getState();
    expect(state.columns).toHaveLength(3);
    expect(state.columns[1]!.id).toBe("c1.5");
  });
});
