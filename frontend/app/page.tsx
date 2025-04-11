import Image from "next/image";
import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function Home() {
  return (
    <div className="flex flex-col min-h-screen">
      <header className="border-b">
        <div className="container mx-auto flex h-16 items-center justify-between px-4">
          <div className="flex items-center gap-2">
            <Image 
              src="/papo-reto-icon.svg" 
              alt="Papo Reto Logo" 
              width={40} 
              height={40}
              className="h-8 w-8"
            />
            <span className="text-xl font-bold">Papo Reto</span>
          </div>
          <nav className="flex items-center gap-4">
            <Link href="/login" className="text-sm font-medium hover:underline">
              Entrar
            </Link>
            <Link href="/register">
              <Button>Criar Conta</Button>
            </Link>
          </nav>
        </div>
      </header>
      <main className="flex-1">
        <section className="container mx-auto px-4 py-12 md:py-24 lg:py-32">
          <div className="grid gap-6 lg:grid-cols-2 lg:gap-12 items-center">
            <div className="space-y-4">
              <h1 className="text-4xl font-bold tracking-tighter sm:text-5xl md:text-6xl">
                Receba mensagens anônimas de forma organizada
              </h1>
              <p className="text-gray-500 md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed dark:text-gray-400">
                Crie grupos personalizados, compartilhe links e receba feedback honesto. 
                Organize suas mensagens anônimas em um só lugar.
              </p>
              <div className="flex flex-col gap-2 min-[400px]:flex-row">
                <Link href="/register">
                  <Button size="lg" className="w-full min-[400px]:w-auto">
                    Começar Grátis
                  </Button>
                </Link>
                <Link href="/about">
                  <Button size="lg" variant="outline" className="w-full min-[400px]:w-auto">
                    Saiba Mais
                  </Button>
                </Link>
              </div>
            </div>
            <div className="flex justify-center">
              <Image
                src="/hero-image.svg"
                alt="Ilustração de mensagens"
                width={1000}
                height={1000}
                className="rounded-lg object-cover"
              />
            </div>
          </div>
        </section>
        <section className="bg-gray-50 dark:bg-gray-900">
          <div className="container mx-auto px-4 py-12 md:py-24 lg:py-32">
            <div className="grid gap-6 lg:grid-cols-3 lg:gap-12">
              <div className="space-y-2">
                <h3 className="text-xl font-bold">Grupos Personalizados</h3>
                <p className="text-gray-500 dark:text-gray-400">
                  Crie grupos diferentes para cada contexto e personalize as configurações.
                </p>
              </div>
              <div className="space-y-2">
                <h3 className="text-xl font-bold">Compartilhamento Fácil</h3>
                <p className="text-gray-500 dark:text-gray-400">
                  Compartilhe links para receber mensagens anônimas em qualquer plataforma.
                </p>
              </div>
              <div className="space-y-2">
                <h3 className="text-xl font-bold">Organização Eficiente</h3>
                <p className="text-gray-500 dark:text-gray-400">
                  Mantenha suas mensagens organizadas, favoritas e filtradas.
                </p>
              </div>
            </div>
          </div>
        </section>
      </main>
      <footer className="border-t">
        <div className="container mx-auto px-4 py-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div className="text-sm text-gray-500 dark:text-gray-400">
            2025 Papo Reto. Todos os direitos reservados.
          </div>
          <nav className="flex gap-4 text-sm text-gray-500 dark:text-gray-400">
            <Link href="/terms" className="hover:underline">
              Termos
            </Link>
            <Link href="/privacy" className="hover:underline">
              Privacidade
            </Link>
            <Link href="/contact" className="hover:underline">
              Contato
            </Link>
          </nav>
        </div>
      </footer>
    </div>
  );
}
