import { MoreHorizontal } from 'lucide-react'
import type { Quest } from '@/apis/types'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { cn } from '@/lib/utils'
import { difficultyClass, difficultyLabel } from '../lib/difficulty'

/**
 * Tabel definisi quest untuk halaman Quests (F2.9). Presentational-ish: terima
 * daftar + callback aksi, tak memanggil hook mutation sendiri.
 */
type QuestTableProps = {
  quests: Quest[]
  onEdit: (quest: Quest) => void
  onArchive: (quest: Quest) => void
}

export function QuestTable({ quests, onEdit, onArchive }: QuestTableProps) {
  return (
    <div className="w-full overflow-x-auto rounded-lg border border-border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Judul</TableHead>
            <TableHead>Kategori</TableHead>
            <TableHead>Difficulty</TableHead>
            <TableHead className="text-right">Poin</TableHead>
            <TableHead>Status</TableHead>
            <TableHead className="w-0 text-right">Aksi</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {quests.map((quest) => (
            <TableRow key={quest.id}>
              <TableCell className="font-medium">{quest.title}</TableCell>
              <TableCell>{quest.category}</TableCell>
              <TableCell>
                <Badge
                  variant="outline"
                  className={cn(
                    'border-transparent',
                    difficultyClass[quest.difficulty],
                  )}
                >
                  {difficultyLabel[quest.difficulty]}
                </Badge>
              </TableCell>
              <TableCell className="text-right tabular-nums">
                {quest.points}
              </TableCell>
              <TableCell>{quest.active ? 'Aktif' : 'Arsip'}</TableCell>
              <TableCell className="text-right">
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="ghost"
                      size="icon-sm"
                      aria-label={`Aksi untuk ${quest.title}`}
                    >
                      <MoreHorizontal />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem onClick={() => onEdit(quest)}>
                      Edit
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={() => onArchive(quest)}>
                      Arsipkan
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
