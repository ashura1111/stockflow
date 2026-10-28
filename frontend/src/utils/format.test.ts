import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { formatCurrency, formatRelativeTime, formatDate } from './format';

describe('Format utilities', () => {
  describe('formatCurrency', () => {
    it('should format number as Indonesian currency', () => {
      const result = formatCurrency(1000000);
      expect(result).toContain('1');
      expect(result).toContain('000');
      expect(result).toContain('000');
    });

    it('should handle zero', () => {
      const result = formatCurrency(0);
      expect(result).toBeTruthy();
    });

    it('should handle decimal values', () => {
      const result = formatCurrency(1500.5);
      expect(result).toBeTruthy();
    });

    it('should have no fraction digits', () => {
      const result = formatCurrency(1000);
      expect(result).not.toContain(',');
    });
  });

  describe('formatRelativeTime', () => {
    let now: Date;

    beforeEach(() => {
      now = new Date('2026-10-28T03:55:00Z');
      vi.useFakeTimers();
      vi.setSystemTime(now);
    });

    afterEach(() => {
      vi.useRealTimers();
    });

    it('should show "just now" for times less than 60 seconds ago', () => {
      const time = new Date(now.getTime() - 30000).toISOString();
      expect(formatRelativeTime(time)).toBe('just now');
    });

    it('should show minutes ago for times less than 1 hour', () => {
      const time = new Date(now.getTime() - 5 * 60000).toISOString();
      expect(formatRelativeTime(time)).toContain('5m ago');
    });

    it('should show hours ago for times less than 1 day', () => {
      const time = new Date(now.getTime() - 3 * 3600000).toISOString();
      expect(formatRelativeTime(time)).toContain('3h ago');
    });

    it('should show days ago for times less than 1 week', () => {
      const time = new Date(now.getTime() - 2 * 86400000).toISOString();
      expect(formatRelativeTime(time)).toContain('2d ago');
    });

    it('should show full date for older times', () => {
      const time = new Date(now.getTime() - 10 * 86400000).toISOString();
      const result = formatRelativeTime(time);
      expect(result).toMatch(/\d/);
      expect(result).not.toContain('ago');
    });
  });

  describe('formatDate', () => {
    it('should format date with time', () => {
      const dateString = '2026-10-28T03:55:00Z';
      const result = formatDate(dateString);
      expect(result).toContain('28');
      expect(result).toContain('2026');
    });

    it('should include both date and time components', () => {
      const dateString = '2026-10-28T15:30:45Z';
      const result = formatDate(dateString);
      expect(result).toMatch(/\d{2}\.\d{2}/);
    });
  });
});
