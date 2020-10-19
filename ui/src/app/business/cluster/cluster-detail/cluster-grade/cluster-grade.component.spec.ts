import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ClusterGradeComponent } from './cluster-grade.component';

describe('ClusterGradeComponent', () => {
  let component: ClusterGradeComponent;
  let fixture: ComponentFixture<ClusterGradeComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ClusterGradeComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(ClusterGradeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
