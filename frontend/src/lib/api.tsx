export type Unit = {
    id: string;
    name: string;
    type: string;
    status: string;
};

export type PaginatedUnits = {
    content: Unit[];
    page: number;
    size: number;
    total: number;
    totalPages: number;
};

const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://127.0.0.1:5000'

export async function fetchUnits(page: number = 0, size: number = 10): Promise<PaginatedUnits> {
    const res = await fetch(`${apiUrl}/api/unit?page=${page}&size=${size}`, {
        cache: "no-store",
    });

    if (!res.ok) {
        throw new Error("Failed to fetch units");
    }

    return res.json().then((data) => data.data);
}
