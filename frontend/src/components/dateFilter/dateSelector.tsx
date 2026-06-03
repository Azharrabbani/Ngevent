import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";

interface Props {
    value: Date | null;
    onChange: (date: Date | null) => void;
}

export default function DateSelector({
    value,
    onChange,
}: Props) {
    return (
        <DatePicker
            selected={value}
            onChange={onChange}
            className="w-full border px-3 py-2 rounded"
            placeholderText="Select date"
        />
    );
}