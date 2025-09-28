import { Unit } from "./Unit";

export type Pagination = {
    page: number;
    size: number;
    total: number;
    totalPages: number;
};

export type PaginatedUnits = {
    content: Unit[];
    pagination: Pagination;
};