import type { PathPoint } from "./publicEventResponse";

export interface RouteTestRequest {
    from_lat: number;
    from_lon: number;
    to_lat: number;
    to_lon: number;
    to_name: string;
}

export interface RouteTestResponse {
    event: string;
    distance: string;
    path: PathPoint[];
    dijkstra_cost: number;
    dijkstra_time_ms: number;
}