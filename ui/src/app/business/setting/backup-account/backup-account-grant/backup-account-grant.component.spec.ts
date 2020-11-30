import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BackupAccountGrantComponent } from './backup-account-grant.component';

describe('BackupAccountGrantComponent', () => {
  let component: BackupAccountGrantComponent;
  let fixture: ComponentFixture<BackupAccountGrantComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BackupAccountGrantComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BackupAccountGrantComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
