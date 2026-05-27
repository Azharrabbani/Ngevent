export interface categoriesResp {
    id: number;
    name: string;
    slug: string;
    created_at: number;
    updated_at: number;
};

export interface categoriesPaginatedResp {
    id: string;
    name: string;
    slug: string;
    total_used: number;
};