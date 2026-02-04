// TODO - Messy

export interface Refill {
  id: string;
  userId: string;
  date: Date;
  weightKg: number;
  bags: number;
  temperature?: number;
  season?: string;
}
