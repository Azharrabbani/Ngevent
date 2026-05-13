import { useEditor, EditorContent } from "@tiptap/react"
import StarterKit from "@tiptap/starter-kit"
import TextAlign from "@tiptap/extension-text-align"
import Underline from "@tiptap/extension-underline"
import {
    FaBold, FaItalic, FaUnderline,
    FaAlignLeft, FaAlignCenter, FaAlignRight,
    FaListUl, FaListOl
} from "react-icons/fa"
import React, { useEffect } from "react"

interface Props {
    value: string;
    onChange: (value: string) => void;
    error?: string;
};

export default function RichTextEditor({ value, onChange, error }: Props) {
    const editor = useEditor({
        extensions: [
            StarterKit,
            Underline,
            TextAlign.configure({
                types: ["heading", "paragraph"]
            }),
        ],
        content: value,
        onUpdate: ({ editor }) => {
            onChange(editor.getHTML())
        }
    })

    if (!editor) return null;

    const ToolbarButton = ({
        onClick,
        active,
        children,
    }: {
        onClick: () => void
        active?: boolean
        children: React.ReactNode
    }) => (
        <button
            type="button"
            onClick={onClick}
            className={`p-2 rounded hover:bg-gray-200 transition ${active ? "bg-gray-300 text-blue-600" : "text-gray-600"
                }`}
        >
            {children}
        </button>
    )

    useEffect(() => {
        if (!editor) return
        if (value === editor.getHTML()) return
        editor.commands.setContent(value)
    }, [value, editor])

    return (
        <div className={`rounded-xl border ${error ? "border-red-500" : "border-gray-300"}`}>
            {/* Toolbar */}
            <div className="flex flex-wrap gap-1 p-2 border-b border-gray-200 bg-gray-50 rounded-t-xl">
                <ToolbarButton
                    onClick={() => editor.chain().focus().toggleBold().run()}
                    active={editor.isActive("bold")}
                >
                    <FaBold size={14} />
                </ToolbarButton>

                <ToolbarButton
                    onClick={() => editor.chain().focus().toggleItalic().run()}
                    active={editor.isActive("italic")}
                >
                    <FaItalic size={14} />
                </ToolbarButton>

                <ToolbarButton
                    onClick={() => editor.chain().focus().toggleUnderline().run()}
                    active={editor.isActive("underline")}
                >
                    <FaUnderline size={14} />
                </ToolbarButton>

                <ToolbarButton
                    onClick={() => editor.chain().focus().setTextAlign("left").run()}
                    active={editor.isActive({ textAlign: "left" })}
                >
                    <FaAlignLeft size={14} />
                </ToolbarButton>

                <ToolbarButton
                    onClick={() => editor.chain().focus().setTextAlign("center").run()}
                    active={editor.isActive({ textAlign: "center" })}
                >
                    <FaAlignCenter size={14} />
                </ToolbarButton>

                <ToolbarButton
                    onClick={() => editor.chain().focus().setTextAlign("right").run()}
                    active={editor.isActive({ textAlign: "right" })}
                >
                    <FaAlignRight size={14} />
                </ToolbarButton>

                <div className="w-px bg-gray-300 mx-1" />

                <ToolbarButton
                    onClick={() => editor.chain().focus().toggleBulletList().run()}
                    active={editor.isActive("bulletList")}
                >
                    <FaListUl size={14} />
                </ToolbarButton>

                <ToolbarButton
                    onClick={() => editor.chain().focus().toggleOrderedList().run()}
                    active={editor.isActive("orderedList")}
                >
                    <FaListOl size={14} />
                </ToolbarButton>

                <div className="w-px bg-gray-300 mx-1" />

                <select
                    className="text-xs border border-gray-300 rounded px-1 bg-white"
                    onChange={(e) => {
                        const level = parseInt(e.target.value)
                        if (level === 0) {
                            editor.chain().focus().setParagraph().run()
                        } else {
                            editor.chain().focus().toggleHeading({ level: level as 1 | 2 | 3 }).run()
                        }
                    }}
                    value={
                        editor.isActive("heading", { level: 1 }) ? 1 :
                            editor.isActive("heading", { level: 2 }) ? 2 :
                                editor.isActive("heading", { level: 3 }) ? 3 : 0
                    }
                >
                    <option value="0">Normal</option>
                    <option value="1">Heading 1</option>
                    <option value="2">Heading 2</option>
                    <option value="3">Heading 3</option>
                </select>
            </div>

            {/* Editor */}
            <EditorContent
                editor={editor}
                className="min-h-[160px] p-4 text-sm outline-none prose max-w-none"
            >
                {error && <p className="text-red-500 text-sm px-4 pb-2">{error}</p>}
            </EditorContent>
        </div>
    )
}