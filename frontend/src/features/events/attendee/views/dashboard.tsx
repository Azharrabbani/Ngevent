import { useGetEventsActive } from "../../hooks/useGeEventsActive"

export default function PublicDashboard() {
    const { data, isLoading } = useGetEventsActive({
        pagination: {
            page: 1,
            limit: 10
        }
    })

    if (isLoading) {
        return <h1>Loading</h1>
    }

    console.log(data);

    return (
        <h1>Public Dashboard</h1>
    )
}