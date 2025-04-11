import Link from "next/link";
import Image from "next/image";
import { RegisterForm } from "@/components/auth/register-form";

export default function LoginPage() {
  return (
    <div className="flex min-h-screen flex-col">
      <div className="flex h-16 items-center border-b px-4">
        <Link href="/" className="flex items-center gap-2">
          <Image 
            src="/papo-reto-icon.svg" 
            alt="Papo Reto Logo" 
            width={40} 
            height={40}
            className="h-8 w-8"
          />
          <span className="text-xl font-bold">Papo Reto</span>
        </Link>
      </div>
      <div className="flex flex-1 items-center justify-center p-4">
        <div className="mx-auto w-full max-w-md space-y-6">
          <RegisterForm />
        </div>
      </div>
    </div>
  );
}
