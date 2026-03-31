import { useCallback, useEffect, useRef, useState } from "react";
import {
  getUploadUrl,
  uploadFileToPresignedUrl,
} from "@/api/card-client";
import type { Editor } from "@tiptap/react";

const MAX_FILE_SIZE = 25 * 1024 * 1024; // 25 MB

export function useDeferredUpload() {
  const pendingUploads = useRef<Map<string, File>>(new Map());
  const [pendingCount, setPendingCount] = useState(0);
  const [isUploading, setIsUploading] = useState(false);

  // Cleanup blob URLs on unmount
  useEffect(() => {
    return () => {
      for (const blobUrl of pendingUploads.current.keys()) {
        URL.revokeObjectURL(blobUrl);
      }
      pendingUploads.current.clear();
    };
  }, []);

  const addFile = useCallback((file: File): string | null => {
    if (file.size > MAX_FILE_SIZE) {
      return null;
    }
    const blobUrl = URL.createObjectURL(file);
    pendingUploads.current.set(blobUrl, file);
    setPendingCount(pendingUploads.current.size);
    return blobUrl;
  }, []);

  const flushUploads = useCallback(
    async (editor: Editor): Promise<void> => {
      const pending = pendingUploads.current;
      if (pending.size === 0) return;

      setIsUploading(true);
      try {
        const entries = Array.from(pending.entries());
        await Promise.all(
          entries.map(async ([blobUrl, file]) => {
            const { upload_url, public_url } = await getUploadUrl({
              filename: file.name,
              content_type: file.type,
            });
            await uploadFileToPresignedUrl(upload_url, file);

            // Replace blob URL with public URL in editor content
            const { doc } = editor.state;
            const tr = editor.state.tr;
            doc.descendants((node, pos) => {
              if (node.type.name === "image" && node.attrs.src === blobUrl) {
                tr.setNodeMarkup(pos, undefined, {
                  ...node.attrs,
                  src: public_url,
                });
              }
            });
            if (tr.docChanged) {
              editor.view.dispatch(tr);
            }

            URL.revokeObjectURL(blobUrl);
            pending.delete(blobUrl);
          }),
        );
      } finally {
        setPendingCount(pendingUploads.current.size);
        setIsUploading(false);
      }
    },
    [],
  );

  return { addFile, flushUploads, pendingCount, isUploading };
}
