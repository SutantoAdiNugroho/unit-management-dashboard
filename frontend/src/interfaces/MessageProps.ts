interface MessageProps {
    open: boolean;
    onClose: () => void;
    title: string;
    message: string;
    isSuccess: boolean;
}