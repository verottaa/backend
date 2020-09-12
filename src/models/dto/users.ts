class UserDto {
  firstName: string;
  lastName: string;
}

export class User extends UserDto {
  id: string;
}

export class CreateUserDto extends UserDto {
}

export class UpdateUserDto extends UserDto {
  id: string;
}
