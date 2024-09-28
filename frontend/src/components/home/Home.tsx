import "./Home.scss"
import { useEffect, useState } from "react";
import { Avatar, Icon, Modal } from "../generic";
import { CardTask } from "./CardTask";
import { TeamModal } from "./TeamModal";

export const Home = () => {
    const [isOpenTeamModal, setIsOpenTeamModal] = useState(false);

    useEffect(() => {
        document.title = "thorac"
    }, []);

    return (
        <div className="Home-container">
            <header>
                <div className="logo">thorac</div>
                <div className="user-bar">
                    <Avatar url="https://i.pinimg.com/originals/9c/2f/56/9c2f562b06908365b4d19d22a50e8f45.jpg" size={40} />
                    <div className="info">
                        <div className="fullName">Ilya Famin</div>
                        <div className="userId">@Kneepy</div>
                    </div>
                </div>
            </header>
            <div className="project-list">
                <div className="project-box">
                    <div
                        className="project"
                        style={{
                            backgroundImage: `url('https://i.pinimg.com/originals/fc/d5/2b/fcd52b427fa342ce3a9e619044cb1b88.jpg')`
                        }}
                    ></div>
                </div>
                <div className="separator"></div>
                <div className="add-new-project">
                    <Icon size={25}>add</Icon>
                </div>
            </div>
            <div className="project-view">
                <div className="name">
                    Test Project
                    {isOpenTeamModal && <TeamModal onClose={() => setIsOpenTeamModal(false)}  users={["sfwefewf"]}/> }
                    <div className="info-project">
                        <button onClick={() => setIsOpenTeamModal(true)} className="team">
                            <Icon size={25}>groups</Icon>
                        </button>
                        <button className="create-task">
                            <Icon size={25}>add</Icon>
                        </button>
                    </div>
                </div>
                <div className="separator"></div>
                <div className="tasks">
                    <CardTask name="Вадим Дебил Тупой" description="Я не шучу" state="В процессе признания"
                              createdAt="сейчас понял" />
                </div>
            </div>
        </div>
    )
}