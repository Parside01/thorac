import React from "react";
import { createPortal } from "react-dom";
import "./Modal.scss"

interface ModalProps {
    children?: React.ReactNode;
    onClose: () => void
}

export const Modal = ({ children, onClose }: ModalProps) => {
    const handleOutsideClick = () => onClose()

    return (
        createPortal(

            <div className="Modal">
                <div onClick={handleOutsideClick} className="modal-background">
                    <div onClick={(e) => e.stopPropagation()} className="modal-container">
                        { children }
                    </div>
                </div>
            </div>

        , document.body)
    )
}