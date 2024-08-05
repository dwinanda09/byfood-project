'use client';

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { Book } from '../types/book';
import { bookAPI } from '../api/books';

interface BookContextType {
  books: Book[];
  loading: boolean;
  error: string | null;
  searchTerm: string;
  selectedYear: string;
  sortBy: string;
  sortOrder: 'asc' | 'desc';
  filteredBooks: Book[];
  setSearchTerm: (term: string) => void;
  setSelectedYear: (year: string) => void;
  setSortBy: (field: string) => void;
  setSortOrder: (order: 'asc' | 'desc') => void;
  addBook: (book: Omit<Book, 'id' | 'created_at' | 'updated_at'>) => Promise<void>;
  updateBook: (id: string, book: Omit<Book, 'id' | 'created_at' | 'updated_at'>) => Promise<void>;
  deleteBook: (id: string) => Promise<void>;
  refreshBooks: () => Promise<void>;
  clearFilters: () => void;
}

const BookContext = createContext<BookContextType | undefined>(undefined);

export const useBooks = () => {
  const context = useContext(BookContext);
  if (context === undefined) {
    throw new Error('useBooks must be used within a BookProvider');
  }
  return context;
};

interface BookProviderProps {
  children: ReactNode;
}

export const BookProvider: React.FC<BookProviderProps> = ({ children }) => {
  const [books, setBooks] = useState<Book[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedYear, setSelectedYear] = useState('');
  const [sortBy, setSortBy] = useState('created_at');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');

  const refreshBooks = async () => {
    try {
      setLoading(true);
      setError(null);
      const fetchedBooks = await bookAPI.getAll();
      setBooks(fetchedBooks);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    refreshBooks();
  }, []);

  const addBook = async (bookData: Omit<Book, 'id' | 'created_at' | 'updated_at'>) => {
    try {
      const newBook = await bookAPI.create(bookData);
      setBooks(prev => [newBook, ...prev]);
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Failed to add book');
    }
  };

  const updateBook = async (id: string, bookData: Omit<Book, 'id' | 'created_at' | 'updated_at'>) => {
    try {
      const updatedBook = await bookAPI.update(id, bookData);
      setBooks(prev => prev.map(book => book.id === id ? updatedBook : book));
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Failed to update book');
    }
  };

  const deleteBook = async (id: string) => {
    try {
      await bookAPI.delete(id);
      setBooks(prev => prev.filter(book => book.id !== id));
    } catch (err) {
      throw new Error(err instanceof Error ? err.message : 'Failed to delete book');
    }
  };

  const clearFilters = () => {
    setSearchTerm('');
    setSelectedYear('');
    setSortBy('created_at');
    setSortOrder('desc');
  };

  // Filter and sort books
  const filteredBooks = React.useMemo(() => {
    let filtered = books;

    // Search filter
    if (searchTerm) {
      filtered = filtered.filter(book => 
        book.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
        book.author.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    // Year filter
    if (selectedYear) {
      filtered = filtered.filter(book => book.year.toString() === selectedYear);
    }

    // Sort
    filtered.sort((a, b) => {
      let aValue: any = a[sortBy as keyof Book];
      let bValue: any = b[sortBy as keyof Book];

      // Handle string comparison for dates
      if (sortBy === 'created_at' || sortBy === 'updated_at') {
        aValue = new Date(aValue || '').getTime();
        bValue = new Date(bValue || '').getTime();
      }

      // Handle string comparison
      if (typeof aValue === 'string' && typeof bValue === 'string') {
        aValue = aValue.toLowerCase();
        bValue = bValue.toLowerCase();
      }

      if (sortOrder === 'asc') {
        return aValue > bValue ? 1 : aValue < bValue ? -1 : 0;
      } else {
        return aValue < bValue ? 1 : aValue > bValue ? -1 : 0;
      }
    });

    return filtered;
  }, [books, searchTerm, selectedYear, sortBy, sortOrder]);

  const value: BookContextType = {
    books,
    loading,
    error,
    searchTerm,
    selectedYear,
    sortBy,
    sortOrder,
    filteredBooks,
    setSearchTerm,
    setSelectedYear,
    setSortBy,
    setSortOrder,
    addBook,
    updateBook,
    deleteBook,
    refreshBooks,
    clearFilters,
  };

  return (
    <BookContext.Provider value={value}>
      {children}
    </BookContext.Provider>
  );
};