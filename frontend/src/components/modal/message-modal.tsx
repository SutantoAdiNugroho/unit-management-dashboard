"use client"

import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
} from "@/components/ui/dialog"
import { CheckCircle, XCircle } from "lucide-react"

export default function MessageModal({ open, onClose, title, message, isSuccess }: MessageProps) {
    return (
        <Dialog open={open} onOpenChange={onClose}>
            <DialogContent>
                <DialogHeader className="flex flex-col items-center justify-center text-center">
                    {isSuccess ? (
                        <CheckCircle className="h-16 w-16 text-green-500 mb-2" />
                    ) : (
                        <XCircle className="h-16 w-16 text-red-500 mb-2" />
                    )}
                    <DialogTitle className="text-xl font-bold">{title}</DialogTitle>
                    <DialogDescription className="text-gray-500 mt-2">
                        {message}
                    </DialogDescription>
                </DialogHeader>
            </DialogContent>
        </Dialog>
    )
}