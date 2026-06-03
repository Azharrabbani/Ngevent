export const defaultPagination = (page: number) => ({
    page: page,
    limit: 8,
});

export const getPage = (
    currentPage: number,
    totalPage: number,
    visibleCount: number = 5
) => {
    let start = currentPage;
    let end = currentPage + visibleCount - 1;

    if (end > totalPage) {
        end = totalPage;
        start = Math.max(end - visibleCount + 1, 1);
    }

    return Array.from({ length: end - start + 1 }, (_, i) => start + i);
};