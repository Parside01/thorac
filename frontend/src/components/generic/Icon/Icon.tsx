import React from "react";

interface IconProps {
    rounded?: boolean
    size?: number | null
    color?: string | null
    children: React.ReactNode
}

export const Icon = ({ rounded = true, size = null, color = null, children }: IconProps) => {
    return (
        <span
            className={ rounded ? `material-symbols-rounded` : `material-symbols-outlined` }
            style={{
                fontSize: size ?? `${size}px`,
                color: color ?? ""
            }}
        >
            { children }
        </span>
    )
}