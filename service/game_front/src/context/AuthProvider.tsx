import { useNavigate } from "react-router-dom";
import { useEffect, useState, JSX, createContext, useContext } from "react";
import { User } from "../models/models";
import { fetchLogout, fetchUser } from "../api/user";

interface AuthContextType {
    user: User | null;
    isLoading: boolean;
    logout: () => void;
    setUser: (user: User | null) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: JSX.Element }) => {
    const navigate = useNavigate();
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(true);

    useEffect(() => {
        const getUser = async () => {
        setIsLoading(true);
        try {
            setUser(await fetchUser());
        } catch (e) {
            setUser(null);
            console.log("Пользователь не авторизован");
        }
        setIsLoading(false);
        };
        
        getUser();
    }, []);

    useEffect(() => {
        if (!isLoading && !user) {
            navigate("/login");
        }
    }, [isLoading, user, navigate]);

    if (isLoading) {
        return <div className="text-center mt-5">Загрузка...</div>;
    }

    return (
        <AuthContext.Provider value={{ user, isLoading, logout: fetchLogout, setUser }}>
        {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
};