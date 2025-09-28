"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { createUnit, updateUnit } from "@/lib/api";
import { Loader2 } from "lucide-react";
import { Label } from "../ui/label";
import { Formik, Form, Field } from "formik";
import { Unit } from "@/types/Unit";

import MessageModal from "./message-modal";

export default function UnitModal({ open, onClose, unit }: { open: boolean; onClose: () => void; unit?: Unit | null }) {
    const [loading, setLoading] = useState<boolean>(false);
    const [showStatusModal, setShowStatusModal] = useState<boolean>(false);

    const [statusTitle, setStatusTitle] = useState<string>("");
    const [statusMessage, setStatusMessage] = useState<string>("");

    const [isSuccess, setIsSuccess] = useState<boolean>(false);

    const handleCloseMessageSuccess = () => {
        setShowStatusModal(false);
        onClose();
    };

    const handleCloseMessageFailed = () => {
        setShowStatusModal(false);
    };

    return (
        <>
            <Formik
                initialValues={{
                    name: unit?.name || "",
                    status: unit?.status || "",
                    type: unit?.type || "",
                }}
                enableReinitialize={true}
                onSubmit={async (values, { setSubmitting, resetForm }) => {

                    setLoading(true);

                    try {
                        const result = unit
                            ? await updateUnit(unit.id, values)
                            : await createUnit(values);

                        if (result && result.success) {
                            setStatusTitle("Success");
                            setStatusMessage((unit) ? "Unit updated" : "Unit created");
                            setIsSuccess(true);
                            resetForm();
                        } else {
                            setStatusTitle("Failed");
                            setStatusMessage(result?.message);
                            setIsSuccess(false);
                        }
                    } catch (err: any) {
                        console.error("Error submit unit:", err);

                        setStatusTitle("Error");
                        setStatusMessage(err.message);
                        setIsSuccess(false);
                    } finally {
                        setLoading(false);
                        setSubmitting(false);
                        setShowStatusModal(true);
                    }
                }}
            >
                {({ values, setFieldValue, resetForm }) => (
                    <Dialog open={open} onOpenChange={() => {
                        resetForm();
                        onClose();
                    }}>
                        <DialogContent>
                            <DialogHeader>
                                <DialogTitle>{unit ? "Edit Unit" : "Create Unit"}</DialogTitle>
                            </DialogHeader>
                            <Form>
                                <div className="grid gap-4 py-4">
                                    <div className="grid grid-cols-4 items-center gap-4">
                                        <Label htmlFor="name" className="text-right">
                                            Name
                                        </Label>
                                        <Field
                                            as={Input}
                                            id="name"
                                            name="name"
                                            className="col-span-3"
                                        />
                                    </div>
                                    <div className="grid grid-cols-4 items-center gap-4">
                                        <Label htmlFor="status" className="text-right">
                                            Status
                                        </Label>
                                        <Select
                                            value={values.status}
                                            onValueChange={(value) => setFieldValue("status", value)}
                                        >
                                            <SelectTrigger className="col-span-3">
                                                <SelectValue placeholder="Select status" />
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectItem value="Available">Available</SelectItem>
                                                <SelectItem value="Occupied">Occupied</SelectItem>
                                                <SelectItem value="Cleaning In Progress">Cleaning In Progress</SelectItem>
                                                <SelectItem value="Maintenance Needed">Maintenance Needed</SelectItem>
                                            </SelectContent>
                                        </Select>
                                    </div>
                                    <div className="grid grid-cols-4 items-center gap-4">
                                        <Label htmlFor="type" className="text-right">
                                            Type
                                        </Label>
                                        <Select
                                            value={values.type}
                                            onValueChange={(value) => setFieldValue("type", value)}
                                        >
                                            <SelectTrigger className="col-span-3">
                                                <SelectValue placeholder="Select type" />
                                            </SelectTrigger>
                                            <SelectContent>
                                                <SelectItem value="capsule">Capsule</SelectItem>
                                                <SelectItem value="cabin">Cabin</SelectItem>
                                            </SelectContent>
                                        </Select>
                                    </div>
                                </div>
                                <div className="flex justify-end mt-4">
                                    <Button
                                        type="submit"
                                        disabled={loading}
                                        className="bg-primary hover:bg-primary/80 text-white"
                                    >
                                        {loading ? (<Loader2 className="h-4 w-4 animate-spin" />) : unit ? "Update" : "Create"}
                                    </Button>
                                </div>
                            </Form>
                        </DialogContent>
                    </Dialog>
                )}
            </Formik>

            <MessageModal
                open={showStatusModal}
                onClose={isSuccess ? handleCloseMessageSuccess : handleCloseMessageFailed}
                title={statusTitle}
                message={statusMessage}
                isSuccess={isSuccess}
            />

        </>
    );
}