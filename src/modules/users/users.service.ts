import { Injectable, NotFoundException } from '@nestjs/common';
import { User } from 'src/models/users';
import { CreateUserDto, UpdateUserDto } from '../../models/dto/users';

import { v4 as uuidv4 } from 'uuid';


@Injectable()
export class UsersService {
  private users: User[] = [];

  public createUser(user: CreateUserDto) {
    const userId = uuidv4();
    const newUser: User = {
      ...user,
      id: userId,
    };
    this.users.push(newUser);
    return userId;
  }

  public getAllUsers() {
    return this.users;
  }

  public getUserById(id: string) {
    const [user] = this.findUser(id);
    return user;
  }

  public updateUser(id: string, updateUser: UpdateUserDto) {
    const [user, index] = this.findUser(id);
    for (const key of Object.keys(updateUser)) {
      user[key] = user[key] === updateUser[key] ? user[key] : updateUser[key];
    }
    this.users[index] = { ...user };
  }

  public deleteUser(id: string) {
    const [, userIndex] = this.findUser(id);
    this.users.slice(userIndex, 1);
  }

  private findUser(id: string): [User, number] {
    const user = this.users.find(u => u.id === id);
    if (!user) {
      throw new NotFoundException('Could not find user');
    }
    const index = this.users.findIndex(u => u.id === user.id);
    return [user, index];
  }
}
