"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

import { Loader2, Plus, Edit, Trash2, ChevronRight, ChevronLeft } from "lucide-react";
import { Unit } from "@/types/Unit";

import { fetchUnits, deleteUnit } from "@/lib/api";

import UnitModal from "@/components/modal/unit-modal";
import DeleteConfirmationModal from "@/components/modal/delete-unit-modal";
import MessageModal from "@/components/modal/message-modal";

export default function UnitPage() {
  const [units, setUnits] = useState<Unit[]>([]);
  const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [deleteModalOpen, setDeleteModalOpen] = useState(false);
  const [selectedUnit, setSelectedUnit] = useState<Unit | null>(null);
  const [searchParams, setSearchParams] = useState<SearchParams>({
    name: "",
    status: "all",
    type: "all",
  });

  const [page, setPage] = useState(1);
  const [size, setSize] = useState(10);
  const [total, setTotal] = useState(0);

  const [showStatusModal, setShowStatusModal] = useState<boolean>(false)
  const [statusTitle, setStatusTitle] = useState<string>("")
  const [statusMessage, setStatusMessage] = useState<string>("")
  const [isSuccess, setIsSuccess] = useState<boolean>(false)

  const loadUnits = async () => {
    setLoading(true);
    try {
      const res = await fetchUnits(
        page,
        size,
        searchParams.name,
        searchParams.status
      );
      setUnits(res.content);
      setTotal(res.pagination.total);
    } catch (err) {
      console.error("Error fetch units:", err)
      handleMessageError(false);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadUnits();
  }, [page, size, searchParams]);

  const handleCreate = () => {
    setSelectedUnit(null);
    setModalOpen(true)
  };

  const handleEdit = (unit: Unit) => {
    setSelectedUnit(unit);
    setModalOpen(true);
  };

  const handleDeleteClick = (unit: Unit) => {
    setSelectedUnit(unit)
    setDeleteModalOpen(true);
  };

  const handleMessageError = (isDelete: boolean) => {
    setStatusTitle("Failed")
    setStatusMessage((isDelete) ? "Failed to delete unit due to unable connect to server" : "Failed to retrieve units due to unable connect to server")
    setIsSuccess(false)
  }

  const handleDeleteConfirm = async () => {
    if (selectedUnit) {
      try {
        const res = await deleteUnit(selectedUnit.id);
        if (!res.success) {
          handleMessageError(true);
        } else {
          setDeleteModalOpen(false);
          loadUnits()
        }
      } catch (err) {
        console.error("Failed to delete unit:", err);
        handleMessageError(true)
      } finally {
        setSelectedUnit(null)
      }
    }
  };

  console.log('selectedUnit: ', selectedUnit);


  const handleCloseModal = () => {
    setSelectedUnit(null);
    setModalOpen(false);
    loadUnits();
  };

  const handleCloseDeleteModal = () => {
    setSelectedUnit(null);
    setDeleteModalOpen(false);
  };

  const handleSearchChange = (key: keyof SearchParams, value: string) => {
    setSearchParams((prev) => ({ ...prev, [key]: value }));
    setPage(1);
  };

  const handleSizeChange = (value: string) => {
    setSize(Number(value))
    setPage(1);
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-semibold mb-8 text-foreground">
        Unit Management
      </h1>

      <div className="flex flex-col md:flex-row gap-4 mb-6 items-end">
        <div className="flex-1">
          <label htmlFor="search-name" className="block text-sm font-medium text-gray-700 mb-2">
            Search by name
          </label>
          <Input
            id="search-name"
            placeholder="Search.."
            value={searchParams.name}
            onChange={(e) => handleSearchChange("name", e.target.value)}
          />
        </div>
        <div className="flex-1">
          <label htmlFor="filter-status" className="block text-sm font-medium text-gray-700 mb-2">
            Filter status
          </label>
          <Select
            value={searchParams.status}
            onValueChange={(value) => handleSearchChange("status", value)}
          >
            <SelectTrigger id="filter-status">
              <SelectValue placeholder="All Status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All</SelectItem>
              <SelectItem value="Available">Available</SelectItem>
              <SelectItem value="Occupied">Occupied</SelectItem>
              <SelectItem value="Cleaning In Progress">Cleaning In Progress</SelectItem>
              <SelectItem value="Maintenance Needed">Maintenance Needed</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <Button onClick={handleCreate} className="w-full md:w-auto text-white">
          Create Unit <Plus className="mr-2 h-4 w-4" />
        </Button>
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Type</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={5} className="text-center">
                  <div className="flex justify-center items-center py-4">
                    <Loader2 className="h-6 w-6 animate-spin" />
                  </div>
                </TableCell>
              </TableRow>
            ) : units.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} className="text-center">
                  No units found
                </TableCell>
              </TableRow>
            ) : (
              units.map((unit) => (
                <TableRow key={unit.id}>
                  <TableCell>{unit.id}</TableCell>
                  <TableCell>{unit.name}</TableCell>
                  <TableCell>{unit.status}</TableCell>
                  <TableCell>
                    <Badge
                      className={
                        unit.type === "capsule"
                          ? "bg-blue-400 text-white"
                          : "bg-purple-400 text-white"
                      }
                    >
                      {unit.type.toUpperCase()}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-right flex justify-end gap-2">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleEdit(unit)}
                    >
                      <Edit className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleDeleteClick(unit)}
                    >
                      <Trash2 className="h-4 w-4 text-red-500" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      <div className="flex flex-col md:flex-row justify-between items-center mt-4">
        <p className="text-sm text-gray-500 mb-2 md:mb-0">
          Total {total} units
        </p>

        <div className="flex flex-col-reverse md:flex-row items-center gap-2">
          <div className="flex items-center gap-2">
            <span className="text-sm text-gray-700">Table size</span>
            <Select value={String(size)} onValueChange={handleSizeChange}>
              <SelectTrigger className="w-[80px]">
                <SelectValue placeholder="10" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="5">5</SelectItem>
                <SelectItem value="10">10</SelectItem>
                <SelectItem value="20">20</SelectItem>
                <SelectItem value="25">25</SelectItem>
                <SelectItem value="50">50</SelectItem>
                <SelectItem value="100">100</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="flex gap-2">
            <Button
              variant="outline"
              onClick={() => setPage((prev) => Math.max(1, prev - 1))}
              disabled={page === 1}
            >
              <ChevronLeft className="h-4 w-4" />
            </Button>
            <Button
              variant="outline"
              onClick={() => setPage((prev) => prev + 1)}
              disabled={page * size >= total}
            >
              <ChevronRight className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </div>

      <UnitModal
        open={modalOpen}
        onClose={handleCloseModal}
        unit={selectedUnit}
      />

      <DeleteConfirmationModal
        open={deleteModalOpen}
        onClose={handleCloseDeleteModal}
        onConfirm={handleDeleteConfirm}
        unitName={selectedUnit?.name || ""}
      />

      <MessageModal
        open={showStatusModal}
        onClose={() => { }}
        title={statusTitle}
        message={statusMessage}
        isSuccess={isSuccess}
      />
    </div>
  );
}