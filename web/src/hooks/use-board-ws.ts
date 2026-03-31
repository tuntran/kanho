import { useEffect, useRef } from "react";
import { useAuthStore } from "@/stores/auth-store";
import { useBoardStore } from "@/stores/board-store";
import type { WSEvent } from "@/types/board";

const RECONNECT_BASE_MS = 500;
const RECONNECT_MAX_MS = 30_000;

export function useBoardWs(boardId: string | null, refetch: () => Promise<void> | void) {
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectAttempts = useRef<number>(0);
  const reconnectTimer = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);

  useEffect(() => {
    if (!boardId) return;

    const connect = () => {
      const token = useAuthStore.getState().accessToken;
      if (!token) return;

      const proto = window.location.protocol === "https:" ? "wss:" : "ws:";
      const url = `${proto}//${window.location.host}/api/v1/boards/${boardId}/ws?token=${token}`;

      const ws = new WebSocket(url);
      wsRef.current = ws;

      ws.onopen = () => {
        reconnectAttempts.current = 0;
        useBoardStore.getState().setWsConnected(true);

        // Subscribe to board
        ws.send(JSON.stringify({ action: "subscribe", board_id: boardId }));
      };

      ws.onmessage = (evt) => {
        try {
          const event = JSON.parse(evt.data as string) as WSEvent;
          if (event.board_id === boardId) {
            useBoardStore.getState().handleServerEvent(event);
          }
        } catch {
          // ignore malformed messages
        }
      };

      ws.onclose = () => {
        useBoardStore.getState().setWsConnected(false);
        wsRef.current = null;
        scheduleReconnect();
      };

      ws.onerror = () => {
        ws.close();
      };
    };

    const scheduleReconnect = () => {
      const delay = Math.min(
        RECONNECT_BASE_MS * Math.pow(2, reconnectAttempts.current),
        RECONNECT_MAX_MS,
      );
      reconnectAttempts.current++;

      reconnectTimer.current = setTimeout(() => {
        connect();
        // Refetch board data on reconnect
        refetch();
      }, delay);
    };

    connect();

    return () => {
      clearTimeout(reconnectTimer.current);
      if (wsRef.current) {
        wsRef.current.onclose = null; // prevent reconnect on intentional close
        wsRef.current.close();
        wsRef.current = null;
      }
      useBoardStore.getState().setWsConnected(false);
    };
  }, [boardId, refetch]);
}
