import { describe, it, expect, beforeEach, vi } from 'vitest';
import { api, ApiError } from './client';

describe('ApiClient', () => {
  beforeEach(() => {
    global.fetch = vi.fn();
  });

  describe('request method', () => {
    it('should make a successful GET request', async () => {
      const mockData = { id: 1, name: 'Test' };
      global.fetch = vi.fn().mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => mockData,
      });

      const result = await api.get('/test');
      expect(result).toEqual(mockData);
      expect(global.fetch).toHaveBeenCalledWith(
        '/api/test',
        expect.objectContaining({ method: 'GET' })
      );
    });

    it('should handle 204 No Content response', async () => {
      global.fetch = vi.fn().mockResolvedValueOnce({
        ok: true,
        status: 204,
      });

      const result = await api.get('/test');
      expect(result).toBeUndefined();
    });

    it('should throw ApiError on failed response', async () => {
      const errorResponse = {
        error: { message: 'Not found', code: 'NOT_FOUND' },
      };
      global.fetch = vi.fn().mockResolvedValueOnce({
        ok: false,
        json: async () => errorResponse,
      });

      await expect(api.get('/test')).rejects.toThrow(ApiError);
    });

    it('should include Content-Type header', async () => {
      global.fetch = vi.fn().mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({}),
      });

      await api.post('/test', { data: 'test' });

      expect(global.fetch).toHaveBeenCalledWith(
        expect.anything(),
        expect.objectContaining({
          headers: { 'Content-Type': 'application/json' },
        })
      );
    });

    it('should serialize body to JSON', async () => {
      global.fetch = vi.fn().mockResolvedValueOnce({
        ok: true,
        status: 200,
        json: async () => ({}),
      });

      const testData = { name: 'Test', value: 42 };
      await api.post('/test', testData);

      expect(global.fetch).toHaveBeenCalledWith(
        expect.anything(),
        expect.objectContaining({
          body: JSON.stringify(testData),
        })
      );
    });
  });

  describe('HTTP methods', () => {
    beforeEach(() => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => ({ success: true }),
      });
    });

    it('should support GET requests', async () => {
      await api.get('/test');
      expect(global.fetch).toHaveBeenCalledWith(
        expect.anything(),
        expect.objectContaining({ method: 'GET' })
      );
    });

    it('should support POST requests', async () => {
      await api.post('/test', {});
      expect(global.fetch).toHaveBeenCalledWith(
        expect.anything(),
        expect.objectContaining({ method: 'POST' })
      );
    });

    it('should support PUT requests', async () => {
      await api.put('/test', {});
      expect(global.fetch).toHaveBeenCalledWith(
        expect.anything(),
        expect.objectContaining({ method: 'PUT' })
      );
    });

    it('should support DELETE requests', async () => {
      await api.delete('/test');
      expect(global.fetch).toHaveBeenCalledWith(
        expect.anything(),
        expect.objectContaining({ method: 'DELETE' })
      );
    });
  });

  describe('API endpoints', () => {
    beforeEach(() => {
      global.fetch = vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: async () => ({}),
      });
    });

    it('should call health endpoint', async () => {
      await api.health();
      expect(global.fetch).toHaveBeenCalledWith('/api/health', expect.anything());
    });

    it('should call dashboard endpoint', async () => {
      await api.dashboard();
      expect(global.fetch).toHaveBeenCalledWith('/api/dashboard', expect.anything());
    });

    it('should build product query with parameters', async () => {
      await api.getProducts('SKU-001', 'active', true, 2, 20);
      const mockFetch = global.fetch as ReturnType<typeof vi.fn>;
      const call = mockFetch.mock.calls[0][0] as string;
      expect(call).toContain('/products?');
      expect(call).toContain('search=SKU-001');
      expect(call).toContain('status=active');
      expect(call).toContain('low_stock=true');
      expect(call).toContain('page=2');
      expect(call).toContain('limit=20');
    });

    it('should build movement query with parameters', async () => {
      await api.getMovements('PROD', 'in', 1, 10);
      const mockFetch = global.fetch as ReturnType<typeof vi.fn>;
      const call = mockFetch.mock.calls[0][0] as string;
      expect(call).toContain('/movements?');
      expect(call).toContain('search=PROD');
      expect(call).toContain('type=in');
      expect(call).toContain('page=1');
      expect(call).toContain('limit=10');
    });
  });

  describe('ApiError', () => {
    it('should create error with correct properties', () => {
      const details = { message: 'Test error', code: 'TEST_ERROR', sku: 'SKU-001' };
      const error = new ApiError('Test message', 'TEST_ERROR', details);

      expect(error.message).toBe('Test message');
      expect(error.code).toBe('TEST_ERROR');
      expect(error.details).toEqual(details);
      expect(error.name).toBe('ApiError');
    });

    it('should be instanceof Error', () => {
      const error = new ApiError('test', 'code', { message: 'test', code: 'code' });
      expect(error).toBeInstanceOf(Error);
    });
  });
});
