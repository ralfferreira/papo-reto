"use client"

import { useEffect, useState } from "react"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import Link from "next/link"
import { groupsService } from "@/lib/api/groups"
// Removed the unused import: import { messagesService } from "@/lib/api/messages"
import { authService } from "@/lib/api/auth"

export default function DashboardPage() {
  const [stats, setStats] = useState({
    activeGroups: 0,
    totalMessages: 0,
    unreadMessages: 0,
    plan: "free"
  })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        // Obter perfil do usuário
        const profileResponse = await authService.getProfile()
        
        // Obter grupos ativos
        const groupsResponse = await groupsService.getGroups(false)
        
        // Calcular estatísticas
        const activeGroups = groupsResponse.data?.groups?.length || 0
        
        setStats({
          activeGroups,
          totalMessages: profileResponse.data?.messageCount || 0,
          unreadMessages: 0, // Seria necessário implementar uma API para contar mensagens não lidas
          plan: profileResponse.data?.plan || "free"
        })
      } catch (error) {
        console.error("Erro ao carregar dados do dashboard:", error)
      } finally {
        setLoading(false)
      }
    }

    fetchDashboardData()
  }, [])

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Dashboard</h2>
          <p className="text-muted-foreground">
            Bem-vindo ao seu painel de controle de mensagens anônimas
          </p>
        </div>
        <Link href="/dashboard/groups/new">
          <Button>Criar Novo Grupo</Button>
        </Link>
      </div>
      
      {loading ? (
        <div className="flex h-40 items-center justify-center">
          <p>Carregando...</p>
        </div>
      ) : (
        <>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium">Grupos Ativos</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.activeGroups}</div>
                <p className="text-xs text-muted-foreground">
                  {stats.plan === "free" ? "Limite: 3 grupos" : "Ilimitado"}
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium">Total de Mensagens</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.totalMessages}</div>
                <p className="text-xs text-muted-foreground">
                  {stats.plan === "free" ? "Limite: 50 por mês" : "Ilimitado"}
                </p>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium">Mensagens Não Lidas</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.unreadMessages}</div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium">Plano Atual</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold capitalize">{stats.plan}</div>
                {stats.plan === "free" && (
                  <p className="text-xs text-muted-foreground">
                    <Link href="/dashboard/settings" className="text-primary hover:underline">
                      Faça upgrade para Premium
                    </Link>
                  </p>
                )}
              </CardContent>
            </Card>
          </div>

          <div className="grid gap-4 md:grid-cols-2">
            <Card className="col-span-1">
              <CardHeader>
                <CardTitle>Grupos Recentes</CardTitle>
                <CardDescription>
                  Seus grupos de mensagens mais recentes
                </CardDescription>
              </CardHeader>
              <CardContent>
                {stats.activeGroups > 0 ? (
                  <div className="space-y-2">
                    <p className="text-sm">Carregue a página de grupos para ver seus grupos ativos.</p>
                  </div>
                ) : (
                  <div className="flex h-20 items-center justify-center rounded-lg border border-dashed">
                    <p className="text-sm text-muted-foreground">
                      Você ainda não tem grupos. Crie seu primeiro grupo!
                    </p>
                  </div>
                )}
              </CardContent>
              <CardFooter>
                <Link href="/dashboard/groups">
                  <Button variant="outline">Ver Todos os Grupos</Button>
                </Link>
              </CardFooter>
            </Card>
            <Card className="col-span-1">
              <CardHeader>
                <CardTitle>Mensagens Recentes</CardTitle>
                <CardDescription>
                  Suas mensagens mais recentes
                </CardDescription>
              </CardHeader>
              <CardContent>
                {stats.totalMessages > 0 ? (
                  <div className="space-y-2">
                    <p className="text-sm">Carregue a página de mensagens para ver suas mensagens recentes.</p>
                  </div>
                ) : (
                  <div className="flex h-20 items-center justify-center rounded-lg border border-dashed">
                    <p className="text-sm text-muted-foreground">
                      Você ainda não recebeu mensagens. Compartilhe seus grupos para começar!
                    </p>
                  </div>
                )}
              </CardContent>
              <CardFooter>
                <Link href="/dashboard/messages">
                  <Button variant="outline">Ver Todas as Mensagens</Button>
                </Link>
              </CardFooter>
            </Card>
          </div>
        </>
      )}
    </div>
  )
}