import { useEffect, useState } from "react"
import { api } from "./api"

export const useAuth = () => {
    const [user, setUser] = useState<any>(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        api.get("/user/me")
        .then(res => {
            setUser(res.data.data)
        })
        .catch(() => {
            setUser(null)
        })
        .finally(() => {
            setLoading(false)
        })
    }, [])

    return {user, loading}
}
