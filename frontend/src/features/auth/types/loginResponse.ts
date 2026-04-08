export interface LoginResponse{
    code: number
    status: string
    message: string
    data: {
      id: string
      email: string
      role: string
      "ngevent-token": string
      "ngevent-ref-token": string
    }
}