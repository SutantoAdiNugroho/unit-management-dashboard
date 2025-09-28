import { PaginatedUnits } from "@/types/Pagination";
import { Unit } from "@/types/Unit";

const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://127.0.0.1:5000/api'

export async function fetchUnits(page: number = 0, size: number = 10, name: string = "", status: string = ""): Promise<PaginatedUnits> {
    if (status == "all") status = "";

    const res = await fetch(`${apiUrl}/unit?page=${page}&size=${size}&name=${name}&status=${status}`, {
        cache: "no-store",
    });

    if (!res.ok) throw new Error("Failed to fetch units")

    return res.json().then((data) => data.data);
}

export async function createUnit(data: Omit<Unit, "id">) {
    const res = await fetch(`${apiUrl}/unit`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
    })

    return res.json()
}

export async function updateUnit(id: string, data: Partial<Unit>) {
    const res = await fetch(`${apiUrl}/unit/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
    })

    return res.json()
}

export async function deleteUnit(id: string) {
    const res = await fetch(`${apiUrl}/unit/${id}`, {
        method: "DELETE",
        headers: { "Content-Type": "application/json" }
    })

    return res.json()
}