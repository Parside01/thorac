import "./TeamModal.scss"
import { Avatar, Modal } from "../../generic";

interface TeamModalProps {
    users: string[]
    onClose: () => void
}

export const TeamModal = ({ users, onClose }: TeamModalProps) => {


    return (
        <Modal onClose={ onClose }>
            <div className="TeamModal">
                <div className="title">Ваша команда</div>
                <div className="users">
                    <div className="user">
                        <Avatar url="https://i.pinimg.com/originals/9c/2f/56/9c2f562b06908365b4d19d22a50e8f45.jpg" size={60} />
                        <div className="userName">Ilya Famin</div>
                    </div>
                </div>
            </div>
        </Modal>
    )
}