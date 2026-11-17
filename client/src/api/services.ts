export interface Service {
  name: string,
  description: string,
  url: string,
  icon: string
}

export interface Category {
  name: string,
  services: Service[]
}

export interface ServiceResponse {
  categories: Category[]
}

export async function getServices(): Promise<ServiceResponse> {
  const response = await fetch('/api/services')

  if (!response.ok) throw new Error('Failed to fetch services');

  return response.json()
}