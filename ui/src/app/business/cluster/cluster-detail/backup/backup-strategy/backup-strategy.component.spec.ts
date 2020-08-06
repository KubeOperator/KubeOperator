import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BackupStrategyComponent } from './backup-strategy.component';

describe('BackupStrategyComponent', () => {
  let component: BackupStrategyComponent;
  let fixture: ComponentFixture<BackupStrategyComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BackupStrategyComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BackupStrategyComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
