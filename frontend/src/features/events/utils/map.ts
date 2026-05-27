import { useEffect } from "react";
import { useMap } from "react-leaflet";

export function ChangeMapView({ position }: { position: [number, number] }) {
    const map = useMap();

    useEffect(() => {
        map.setView(position, 13);
    }, [position, map]);

    return null;
}