import { Book } from '../types/book';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export const bookAPI = {
  async getAll(): Promise<Book[]> {
    const response = await fetch(`${API_URL}/books`);
    if (!response.ok) {
      throw new Error('Failed to fetch books');
    }
    return response.json();
  },

  async getById(id: string): Promise<Book> {
    const response = await fetch(`${API_URL}/books/${id}`);
    if (!response.ok) {
      throw new Error('Failed to fetch book');
    }
    return response.json();
  },

  async create(book: Omit<Book, 'id' | 'created_at' | 'updated_at'>): Promise<Book> {
    const response = await fetch(`${API_URL}/books`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(book),
    });
    if (!response.ok) {
      throw new Error('Failed to create book');
    }
    return response.json();
  },

  async update(id: string, book: Omit<Book, 'id' | 'created_at' | 'updated_at'>): Promise<Book> {
    const response = await fetch(`${API_URL}/books/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(book),
    });
    if (!response.ok) {
      throw new Error('Failed to update book');
    }
    return response.json();
  },

  async delete(id: string): Promise<void> {
    const response = await fetch(`${API_URL}/books/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      throw new Error('Failed to delete book');
    }
  },
};