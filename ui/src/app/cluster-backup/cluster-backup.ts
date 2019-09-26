import {BackupStorage} from '../setting/backup-storage-setting/backup-storage';

export class ClusterBackup {
  id: string;
  name: string;
  folder: string;
  size: number;
  date_created: string;
  project_id: string;
  backup_storage: BackupStorage;
}
