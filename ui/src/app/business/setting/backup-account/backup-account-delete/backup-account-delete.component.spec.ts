import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BackupAccountDeleteComponent } from './backup-account-delete.component';

describe('BackupAccountDeleteComponent', () => {
  let component: BackupAccountDeleteComponent;
  let fixture: ComponentFixture<BackupAccountDeleteComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BackupAccountDeleteComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BackupAccountDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
