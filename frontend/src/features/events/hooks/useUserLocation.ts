import { useState, useEffect } from "react";

interface GeoState {
    lat: number | undefined;
    lon: number | undefined;
    loading: boolean;
    denied: boolean;
}

export const useUserLocation = (): GeoState => {
    const [state, setState] = useState<GeoState>({
        lat: undefined,
        lon: undefined,
        loading: true,
        denied: false,
    });

    useEffect(() => {
        if (!navigator.geolocation) {
            setState((s) => ({ ...s, loading: false, denied: true }));
            return;
        }

        navigator.geolocation.getCurrentPosition(
            (pos) => {
                setState({
                    lat: pos.coords.latitude,
                    lon: pos.coords.longitude,
                    loading: false,
                    denied: false,
                });
            },
            () => {
                setState({ lat: undefined, lon: undefined, loading: false, denied: true });
            }
        );
    }, []);

    return state;
};