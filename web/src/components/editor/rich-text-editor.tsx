import { useCallback, useRef } from "react";
import { useEditor, EditorContent } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import Image from "@tiptap/extension-image";
import Placeholder from "@tiptap/extension-placeholder";
import Link from "@tiptap/extension-link";
import { EditorToolbar } from "./editor-toolbar";
import { useDeferredUpload } from "@/hooks/use-deferred-upload";

interface RichTextEditorProps {
  initialContent: string;
  onSave: (content: string) => void;
  placeholder?: string;
}

export function RichTextEditor({
  initialContent,
  onSave,
  placeholder = "Add description...",
}: RichTextEditorProps) {
  const { addFile, flushUploads, pendingCount, isUploading } =
    useDeferredUpload();
  const fileInputRef = useRef<HTMLInputElement>(null);

  const editor = useEditor({
    extensions: [
      StarterKit,
      Image.configure({ inline: false, allowBase64: false }),
      Placeholder.configure({ placeholder }),
      Link.configure({ openOnClick: false }),
    ],
    content: initialContent || "",
    editorProps: {
      attributes: {
        class:
          "prose prose-invert prose-sm max-w-none px-3 py-2 min-h-[120px] focus:outline-none text-text-1",
      },
      handlePaste: (_view, event) => {
        const items = event.clipboardData?.items;
        if (!items) return false;
        for (const item of items) {
          if (item.type.startsWith("image/")) {
            event.preventDefault();
            const file = item.getAsFile();
            if (!file) continue;
            const blobUrl = addFile(file);
            if (blobUrl) {
              editor?.chain().focus().setImage({ src: blobUrl }).run();
            }
            return true;
          }
        }
        return false;
      },
      handleDrop: (_view, event) => {
        const files = event.dataTransfer?.files;
        if (!files?.length) return false;
        for (const file of files) {
          if (file.type.startsWith("image/")) {
            event.preventDefault();
            const blobUrl = addFile(file);
            if (blobUrl) {
              editor?.chain().focus().setImage({ src: blobUrl }).run();
            }
            return true;
          }
        }
        return false;
      },
    },
    onBlur: async () => {
      if (!editor) return;
      await flushUploads(editor);
      onSave(editor.getHTML());
    },
  });

  const handleImageUpload = useCallback(() => {
    fileInputRef.current?.click();
  }, []);

  const handleFileChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const file = e.target.files?.[0];
      if (!file || !editor) return;
      if (!file.type.startsWith("image/")) return;
      const blobUrl = addFile(file);
      if (blobUrl) {
        editor.chain().focus().setImage({ src: blobUrl }).run();
      }
      e.target.value = "";
    },
    [editor, addFile],
  );

  return (
    <div className="overflow-hidden rounded-lg border border-border bg-surface transition-colors focus-within:border-border-focus">
      <EditorToolbar editor={editor} onImageUpload={handleImageUpload} />
      <EditorContent editor={editor} />
      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        className="hidden"
        onChange={handleFileChange}
      />
      {(pendingCount > 0 || isUploading) && (
        <div className="border-t border-border px-3 py-1.5 text-xs text-text-2">
          {isUploading
            ? "Uploading images..."
            : `${pendingCount} image(s) pending upload`}
        </div>
      )}
    </div>
  );
}
