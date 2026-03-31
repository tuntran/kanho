import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useDeferredUpload } from "./use-deferred-upload";

vi.mock("@/api/card-client", () => ({
  getUploadUrl: vi.fn().mockResolvedValue({
    upload_url: "https://s3.example.com/upload",
    public_url: "https://cdn.example.com/image.png",
  }),
  uploadFileToPresignedUrl: vi.fn().mockResolvedValue(undefined),
}));

// Mock URL.createObjectURL / revokeObjectURL
const mockCreateObjectURL = vi.fn((blob: Blob) => `blob:mock-${blob.size}`);
const mockRevokeObjectURL = vi.fn();
Object.defineProperty(globalThis, "URL", {
  value: {
    ...globalThis.URL,
    createObjectURL: mockCreateObjectURL,
    revokeObjectURL: mockRevokeObjectURL,
  },
  writable: true,
});

describe("useDeferredUpload", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("returns initial state with zero pending", () => {
    const { result } = renderHook(() => useDeferredUpload());

    expect(result.current.pendingCount).toBe(0);
    expect(result.current.isUploading).toBe(false);
    expect(typeof result.current.addFile).toBe("function");
    expect(typeof result.current.flushUploads).toBe("function");
  });

  it("tracks blob URLs added via addFile", () => {
    const { result } = renderHook(() => useDeferredUpload());
    const file = new File(["hello"], "test.png", { type: "image/png" });

    let blobUrl: string | null = null;
    act(() => {
      blobUrl = result.current.addFile(file);
    });

    expect(blobUrl).toBeTruthy();
    expect(blobUrl).toContain("blob:");
    expect(result.current.pendingCount).toBe(1);
  });

  it("rejects files exceeding 25 MB", () => {
    const { result } = renderHook(() => useDeferredUpload());
    // Create a file just over 25 MB
    const bigContent = new ArrayBuffer(26 * 1024 * 1024);
    const file = new File([bigContent], "huge.bin", { type: "application/octet-stream" });

    let blobUrl: string | null = null;
    act(() => {
      blobUrl = result.current.addFile(file);
    });

    expect(blobUrl).toBeNull();
    expect(result.current.pendingCount).toBe(0);
  });

  it("revokes blob URLs on unmount", () => {
    const { result, unmount } = renderHook(() => useDeferredUpload());
    const file = new File(["data"], "img.png", { type: "image/png" });

    act(() => {
      result.current.addFile(file);
    });

    unmount();
    expect(mockRevokeObjectURL).toHaveBeenCalled();
  });
});
