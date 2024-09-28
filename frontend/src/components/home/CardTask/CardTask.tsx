import { Avatar } from "../../generic";
import "./CardTask.scss"

interface CardTaskProps {
    name: string
    description: string
    executors?: string[]
    state: string
    stateColor?: string
    createdAt: string
    completedAt?: string | null
}

// кароче если есть completedAt и статус завершена или типа того то мы вместо времени создания задачи выводим время завершения
export const CardTask = ({ name, description, executors = [], state, stateColor = "#000", createdAt, completedAt = null }: CardTaskProps) => {
    return (
        <div className="CardTask">
            <div className="title">{ name }</div>
            <div className="metadata">
                <div style={{ backgroundColor: stateColor }} className="status">{ state }</div>
                <div className="executors">
                    {
                        executors.map((item, i) => {
                            return (
                                <div
                                    className="executor"
                                    style={{
                                        transform: `translateX(${50 * i}%)`
                                    }}
                                >
                                    <Avatar
                                        size={25}
                                        url="https://i.pinimg.com/originals/b2/d1/c9/b2d1c98dfe0b7585107a2ba8a646c0c6.jpg"
                                    />
                                </div>
                            )
                        })
                    }
                </div>
                <div className="time">{ createdAt }</div>
            </div>
            <div className="separator"></div>
            <div className="description">{ description }</div>
        </div>
    )
}