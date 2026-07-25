import type { ApiError, CalculationRequest, CalculationResponse, OperationManifest } from './types';

const API_BASE = '/api/v1';

const parseErrorBody = async (response: Response): Promise<ApiError> => {
  try {
    return await response.json() as ApiError;
  } catch {
    return { code: 'service_unavailable', message: 'The calculation service is temporarily unavailable.' };
  }
};

export const getOperations = async (): Promise<OperationManifest> => {
  let response: Response;
  try {
    response = await fetch(`${API_BASE}/operations`);
  } catch {
    throw { code: 'manifest_unavailable', message: 'Calculator capabilities are unavailable.' } as ApiError;
  }
  if (!response.ok) {
    throw await parseErrorBody(response);
  }
  return response.json() as Promise<OperationManifest>;
};

export const calculate = async (request: CalculationRequest): Promise<CalculationResponse> => {
  let response: Response;
  try {
    response = await fetch(`${API_BASE}/calculations`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request),
    });
  } catch {
    throw { code: 'service_unavailable', message: 'The calculation service is temporarily unavailable.' } as ApiError;
  }
  if (!response.ok) {
    throw await parseErrorBody(response);
  }
  return response.json() as Promise<CalculationResponse>;
};
