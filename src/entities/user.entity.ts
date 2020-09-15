import { Column, Entity, PrimaryColumn } from 'typeorm';

@Entity()
export class UserEntity {
  @PrimaryColumn({ unique: true })
  id: string;

  @Column({ unique: false })
  firstName: string;

  @Column({ unique: false })
  lastName: string;
}
