import client from "../client";

export interface SetupStatusResponse {
  data: {
    completed: boolean;
  };
}

export interface InitUserRequest {
  username: string;
  password: string;
  roles?: string[];
}

export interface InitUserResponse {
  data: {
    id: string;
    username: string;
    roles: string[];
    isValid: boolean;
    createdAt: string;
    updatedAt: string;
  };
}

export const setupService = {
  async getSetupStatus(): Promise<SetupStatusResponse> {
    const response = await client.get<SetupStatusResponse>("/setup/status");
    return response.data;
  },

  async initUser(userData: InitUserRequest): Promise<InitUserResponse> {
    const response = await client.post<InitUserResponse>(
      "/init-user",
      userData
    );
    return response.data;
  },
};
