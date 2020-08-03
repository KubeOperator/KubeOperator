import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BackupAccountListComponent } from './backup-account-list.component';

describe('BackupAccountListComponent', () => {
  let component: BackupAccountListComponent;
  let fixture: ComponentFixture<BackupAccountListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BackupAccountListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BackupAccountListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
