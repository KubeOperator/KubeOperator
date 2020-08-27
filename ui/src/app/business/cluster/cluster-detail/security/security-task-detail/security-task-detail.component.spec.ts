import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SecurityTaskDetailComponent } from './security-task-detail.component';

describe('SecurityTaskDetailComponent', () => {
  let component: SecurityTaskDetailComponent;
  let fixture: ComponentFixture<SecurityTaskDetailComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SecurityTaskDetailComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SecurityTaskDetailComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
