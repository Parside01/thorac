interface AvatarProps {
    /**
     * Размер Аватара
     */
    size?: number

    /**
     * То насколько сильно закруглены края аваатра
     * По дефолту он круглый
     */
    round?: number

    url: string
}

export const Avatar = ({ size, round, url }: AvatarProps) => {
    return (
        <div
            className="Avatar"
            style={{
                width: `${size}px`,
                height: `${size}px`,
                borderRadius: round ? `${round}px` : "50%",
                backgroundImage: `url('${url}')`,
                backgroundPosition: "center",
                backgroundSize: "cover",
            }}
        ></div>
    )
}