import { IoIosSearch } from "react-icons/io";
import Input from "../../../../components/input";
import { IoFilter } from "react-icons/io5";
import { IoWarningOutline } from "react-icons/io5";

interface Props {
    title: string;
    description: string;
    pendingPage: boolean;

}

export default function PageHeader({ title, description, pendingPage }: Props) {
    return (
        <div className="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between w-full">
            <div className="p-8 space-y-2">
                <h1 className="text-3xl md:text-4xl font-bold">
                    {title}
                </h1>
                <p className="text-sm md:text-base text-[#424654] max-w-2xl">
                    {description}
                </p>
            </div>

            {pendingPage ? (
                <div className="flex bg-[#E2E7FF] rounded-xl px-5 py-3 items-center gap-3 w-full lg:w-auto">
                    <IoWarningOutline className="text-yellow-800" size={23} />
                    <h2 className="font-semibold">12 Awaiting Review</h2>
                </div>
            ) :
                <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 w-full lg:w-auto">
                    <Input
                        className="bg-white border rounded-lg w-full sm:w-72 md:w-80"
                        placeholder="Search events"
                        leftIcon={<IoIosSearch size={24} />}
                    />

                    <button className="flex items-center justify-center gap-2 px-4 py-2 rounded-lg border border-blue-500 text-blue-500 hover:bg-blue-50 transition-all duration-200">
                        <IoFilter size={23} />
                        <span className="text-sm md:text-base font-medium">
                            Filter
                        </span>
                    </button>
                </div>
            }

        </div>
    )
}