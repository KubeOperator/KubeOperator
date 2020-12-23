import { ComponentFixture, TestBed } from '@angular/core/testing';

import { IpCreateComponent } from './ip-create.component';

describe('IpCreateComponent', () => {
  let component: IpCreateComponent;
  let fixture: ComponentFixture<IpCreateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ IpCreateComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(IpCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
