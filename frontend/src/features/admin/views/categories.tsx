import { useEffect, useState } from "react";
import { useListCategoriesPaginated } from "../../categories/hooks/useListCategoriesPaginated";
import CategoriesList from "../components/categories/categoriesList";
import AdminSidebar from "../components/sideBar";
import CategoriesHeader from "../components/categories/header";

export default function CategoriesView() {
    const [currentPage, setCurrentPage] = useState(1);
    const [search, setSearch] = useState<string | undefined>(undefined);

    const { data, isLoading } = useListCategoriesPaginated({
        name: search,
        pagination: {
            limit: 4,
            page: currentPage,
        }
    })

    const totalPage = data?.total_pages ?? 1;

    useEffect(() => {
        const delay = setTimeout(() => {
            setCurrentPage(1);
        }, 500);

        return () => clearTimeout(delay);
    }, [search]);

    return (
        <AdminSidebar>
            <>
                <CategoriesHeader search={search} setSearch={setSearch} />

                <CategoriesList
                    data={data}
                    isLoading={isLoading}
                    search={search}
                    setSearch={setSearch}
                    currentPage={currentPage}
                    totalPage={totalPage}
                    setCurrentPage={setCurrentPage}
                />
            </>
        </AdminSidebar>
    )
}