import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'byFood Library Management System',
  description: 'A production-ready, full-stack CRUD application for managing a library of books.',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}