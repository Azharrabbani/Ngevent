
import { useQuery } from "@tanstack/react-query";

import type { RouteTestRequest } from "../types/routeTest";
import { getRouteTestApi } from "../api/eventsApi";

export const useRouteTest = (params: RouteTestRequest | null) => {
    return useQuery({
        queryKey: ["route-test", params],
        queryFn: () => getRouteTestApi(params as RouteTestRequest),
        enabled: false,
        retry: false,
    });
};